from TTS.tts.configs.xtts_config import XttsConfig
from TTS.tts.models.xtts import Xtts
import os
from scipy.io.wavfile import write
from io import BytesIO


WAV_FILE_PATHS = {
    "female_001": "female-0-100.wav",
    "female_002": "female-0-100.wav",
    "female_003": "female-0-100.wav",
    # 添加更多键值对，每个键都是一个字符串，对应一个WAV文件的路径
}


class ClonerManager:
    def __init__(self,model):
        self.cloners = {}
        self.model = model

    def get_cloner(self,audio_id):
        if audio_id in self.cloners:
            return self.cloners[audio_id]
        else:
            clone = Cloner(self.model, WAV_FILE_PATHS[audio_id])
            self.cloners[audio_id] = clone
            return clone

SAMPLE_RATE = 24_000

class Cloner:
    def __init__(self, model, audio_path):
        self.model = model
        self.gpt_cond_latent, self.diffusion_cond_latents, self.speaker_embedding = self.model.get_conditioning_latents(audio_path)


    def text_to_speech(self, text, language,**kwargs):
       
        outputs = self.model.inference(
                    text,
                    language,
                    self.gpt_cond_latent,
                    self.speaker_embedding,
                    self.diffusion_cond_latents,
                    **kwargs,
                )
        
        # 合成语音
        # outputs = self.synthesize(text, speaker_wav, gpt_cond_len, language)

        audio_stream = BytesIO()
        write(audio_stream, rate=SAMPLE_RATE,data= outputs['wav'])
        
        # Get audio data bytes
        return audio_stream.getvalue()
 # 保存语音到文件
        # torch.save(outputs['wav'], audio_path)

class TTSModel:
    def __init__(self, model_path):
        self.config = XttsConfig()
        self.config.load_json(os.path.join(model_path, 'config.json'))
        self.model = Xtts.init_from_config(self.config)
        self.model.load_checkpoint(self.config, checkpoint_dir=model_path, eval=True, use_deepspeed=True)
        self.model.cuda()
 
    def get_conditioning_latents(
            self,
            audio_path,
            gpt_cond_len=3,
        ): 
        return self.model.get_conditioning_latents(audio_path,gpt_cond_len)
    
    def inference(
        self,
        text,
        language,
        gpt_cond_latent,
        speaker_embedding,
        diffusion_conditioning,
        # GPT inference
        temperature=0.65,
        length_penalty=1,
        repetition_penalty=2.0,
        top_k=50,
        top_p=0.85,
        do_sample=True,
        # Decoder inference
        decoder_iterations=100,
        cond_free=True,
        cond_free_k=2,
        diffusion_temperature=1.0,
        decoder_sampler="ddim",
        **kwargs,
    ):

        # settings = {
        #             "temperature": self.config.temperature,
        #             "length_penalty": self.config.length_penalty,
        #             "repetition_penalty": self.config.repetition_penalty,
        #             "top_k": self.config.top_k,
        #             "top_p": self.config.top_p,
        #             "cond_free_k": self.config.cond_free_k,
        #             "diffusion_temperature": self.config.diffusion_temperature,
        #             "decoder_iterations": self.config.decoder_iterations,
        #             "decoder_sampler": self.config.decoder_sampler,
        #         }
        # settings.update(kwargs)  # allow overriding of preset settings with kwargs
        return self.model.inference(
            text,
            language,
            gpt_cond_latent,
            speaker_embedding,
            diffusion_conditioning,
            # temperature=temperature,
            # length_penalty=length_penalty,
            # repetition_penalty=repetition_penalty,
            # top_k=top_k,
            # top_p=top_p,
            # do_sample=do_sample,
            # decoder_iterations=decoder_iterations,
            # cond_free=cond_free,
            # cond_free_k=cond_free_k,
            # diffusion_temperature=diffusion_temperature,
            # decoder_sampler=decoder_sampler,
            **kwargs,
        )

